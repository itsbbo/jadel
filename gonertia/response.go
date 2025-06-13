package gonertia

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"maps"
	"net/http"
	"strings"
	"sync"
)

// TemplateData are data that will be available in the root template.
type TemplateData map[string]any

// TemplateFuncs are functions that will be available in the root template.
type TemplateFuncs map[string]any

// Props are the data that will be transferred
// and will be available in the front-end component.
type Props map[string]any

// OptionalProp is a property that will evaluate when needed.
//
// https://inertiajs.com/partial-reloads
type OptionalProp struct {
	ignoresFirstLoad
	Value any
}

func (p OptionalProp) Prop() any {
	return p.Value
}

func Optional(value any) OptionalProp {
	return OptionalProp{Value: value}
}

var _ ignoreFirstLoad = OptionalProp{}

type ignoreFirstLoad interface {
	shouldIgnoreFirstLoad() bool
}

type ignoresFirstLoad struct{}

func (i ignoresFirstLoad) shouldIgnoreFirstLoad() bool { return true }

// Deprecated: use OptionalProp.
type LazyProp = OptionalProp

// Deprecated: use Optional.
func Lazy(value any) LazyProp {
	return LazyProp{Value: value}
}

// DeferProp is a property that will evaluate after page load.
//
// https://v2.inertiajs.com/deferred-props
type DeferProp struct {
	ignoresFirstLoad
	mergesProps
	Value any
	Group string
}

func (p DeferProp) Prop() any {
	return p.Value
}

func (p DeferProp) Merge() DeferProp {
	p.merge = true
	return p
}

func Defer(value any, group ...string) DeferProp {
	return DeferProp{
		Value: value,
		Group: firstOr(group, "default"),
	}
}

var _ ignoreFirstLoad = DeferProp{}

var _ mergeable = DeferProp{}

// AlwaysProp is a property that will always evaluated.
//
// https://inertiajs.com/partial-reloads
type AlwaysProp struct {
	Value any
}

func (p AlwaysProp) Prop() any {
	return p.Value
}

func Always(value any) AlwaysProp {
	return AlwaysProp{Value: value}
}

// MergeProps is a property, which items will be merged instead of overwrite.
//
// https://v2.inertiajs.com/merging-props
type MergeProps struct {
	mergesProps
	Value any
}

func (p MergeProps) Prop() any {
	return p.Value
}

func (p MergeProps) Merge() MergeProps {
	p.merge = true
	return p
}

func Merge(value any) MergeProps {
	return MergeProps{
		Value:       value,
		mergesProps: mergesProps{merge: true},
	}
}

var _ mergeable = MergeProps{}

type mergeable interface {
	shouldMerge() bool
}

type mergesProps struct {
	merge bool
}

func (p mergesProps) shouldMerge() bool {
	return p.merge
}

// Proper is an interface for custom type, which provides property, that will be resolved.
type Proper interface {
	Prop() any
}

// ProperWithContext is an interface for custom type, which provides property,
// that will be resolved with context passing.
type ProperWithContext interface {
	PropWithContext(_ context.Context) any
}

// TryProper is an interface for custom type, which provides property and error, that will be resolved.
type TryProper interface {
	TryProp() (any, error)
}

// TryProperWithContext is an interface for custom type, which provides property and error,
// that will be resolved with context passing.
type TryProperWithContext interface {
	TryPropWithContext(_ context.Context) (any, error)
}

// ValidationErrors are messages, that will be stored in the "errors" prop.
type ValidationErrors map[string]any

// Location creates redirect response.
//
// If request was made by Inertia - sets status to 409 and url will be in "X-Inertia-Location" header.
// Otherwise, it will do an HTTP redirect with specified status (default is 302 for GET, 303 for POST/PUT/PATCH).
func (i *Inertia) Location(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	i.flashContext(r.Context())

	if IsInertiaRequest(r) {
		setInertiaLocationInResponse(w, url)
		deleteInertiaInResponse(w)
		deleteVaryInResponse(w)
		setResponseStatus(w, http.StatusConflict)
		return
	}

	redirectResponse(w, r, url, status...)
}

// Back creates plain redirect response to the previous url.
func (i *Inertia) Back(w http.ResponseWriter, r *http.Request, status ...int) {
	i.Redirect(w, r, backURL(r), status...)
}

func backURL(r *http.Request) string {
	// At the moment, it based only on the "Referer" HTTP header.
	return refererFromRequest(r)
}

// Redirect creates plain redirect response.
func (i *Inertia) Redirect(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	i.flashContext(r.Context())

	redirectResponse(w, r, url, status...)
}

func (i *Inertia) flashContext(ctx context.Context) {
	i.flashValidationErrorsFromContext(ctx)
	i.flashClearHistoryFromContext(ctx)
}

func (i *Inertia) flashValidationErrorsFromContext(ctx context.Context) {
	if i.flash == nil {
		return
	}

	validationErrors := ValidationErrorsFromContext(ctx)
	if len(validationErrors) == 0 {
		return
	}

	err := i.flash.FlashErrors(ctx, validationErrors)
	if err != nil {
		i.logger.Printf("cannot flash validation errors: %s", err)
	}
}

func (i *Inertia) flashClearHistoryFromContext(ctx context.Context) {
	if i.flash == nil {
		return
	}

	clearHistory := ClearHistoryFromContext(ctx)
	if !clearHistory {
		return
	}

	err := i.flash.FlashClearHistory(ctx)
	if err != nil {
		i.logger.Printf("cannot flash clear history: %s", err)
	}
}

// Render returns response with Inertia data.
//
// If request was made by Inertia - it will return data in JSON format.
// Otherwise, it will return HTML with root template.
//
// If SSR is enabled, pre-renders JavaScript and return HTML (https://inertiajs.com/server-side-rendering).
func (i *Inertia) Render(w http.ResponseWriter, r *http.Request, component string, props ...Props) (err error) {
	p, err := i.buildPage(r, component, firstOr(props, nil))
	if err != nil {
		return fmt.Errorf("build page: %w", err)
	}

	if IsInertiaRequest(r) {
		if err = i.doInertiaResponse(w, p); err != nil {
			return fmt.Errorf("inertia response: %w", err)
		}
		return
	}

	if err = i.doHTMLResponse(w, r, p); err != nil {
		return fmt.Errorf("html response: %w", err)
	}

	return nil
}

type page struct {
	Component      string              `json:"component"`
	Props          Props               `json:"props"`
	URL            string              `json:"url"`
	Version        string              `json:"version"`
	EncryptHistory bool                `json:"encryptHistory"`
	ClearHistory   bool                `json:"clearHistory"`
	DeferredProps  map[string][]string `json:"deferredProps,omitempty"`
	MergeProps     []string            `json:"mergeProps,omitempty"`
}

func (i *Inertia) buildPage(r *http.Request, component string, props Props) (*page, error) {
	props = i.collectProps(r, props)

	deferredProps := i.resolveDeferredProps(r, component, props)
	mergeProps := resolveMergeProps(r, props)

	props, err := i.resolveProps(r, component, props)
	if err != nil {
		return nil, fmt.Errorf("resolve props: %w", err)
	}

	return &page{
		Component:      component,
		Props:          props,
		URL:            r.RequestURI,
		Version:        i.version,
		EncryptHistory: i.resolveEncryptHistory(r.Context()),
		ClearHistory:   ClearHistoryFromContext(r.Context()),
		DeferredProps:  deferredProps,
		MergeProps:     mergeProps,
	}, nil
}

func (i *Inertia) resolveDeferredProps(r *http.Request, component string, props Props) map[string][]string {
	if isPartial(r, component) {
		return nil
	}

	keysByGroups := make(map[string][]string)

	for key, val := range props {
		if dp, ok := val.(DeferProp); ok {
			keysByGroups[dp.Group] = append(keysByGroups[dp.Group], key)
		}
	}

	return keysByGroups
}

func (i *Inertia) collectProps(r *http.Request, props Props) Props {
	propsCtx := PropsFromContext(r.Context())
	result := make(Props, len(propsCtx)+len(i.sharedProps)+len(props)+1)

	// Add validation errors from context to the result.
	{
		result["errors"] = AlwaysProp{ValidationErrorsFromContext(r.Context())}
	}

	// Add shared props to the result.
	maps.Copy(result, i.sharedProps)

	// Add props from context to the result.
	maps.Copy(result, propsCtx)

	// Add passed props to the result.
	maps.Copy(result, props)

	return result
}

func resolveMergeProps(r *http.Request, props Props) []string {
	resetProps := setOf(resetFromRequest(r))

	var mergeProps []string
	for key, val := range props {
		if _, ok := resetProps[key]; ok {
			continue
		}

		if m, ok := val.(mergeable); ok && m.shouldMerge() {
			mergeProps = append(mergeProps, key)
		}
	}

	return mergeProps
}

//nolint:gocognit
func (i *Inertia) resolveProps(r *http.Request, component string, props Props) (Props, error) {
	// Partial reloads only work for visits made to the same page component.
	//
	// https://inertiajs.com/partial-reloads
	if isPartial(r, component) {
		// Only (include keys) and except (exclude keys) logic.
		only, except := getOnlyAndExcept(r)

		if len(only) > 0 {
			for key, val := range props {
				if _, ok := only[key]; ok {
					continue
				}
				if _, ok := val.(AlwaysProp); ok {
					continue
				}

				delete(props, key)
			}
		}
		for key := range except {
			if _, ok := props[key].(AlwaysProp); ok {
				continue
			}

			delete(props, key)
		}
	} else {
		// Props with ignoreFirstLoad should not be included.
		for key, val := range props {
			if ifl, ok := val.(ignoreFirstLoad); ok && ifl.shouldIgnoreFirstLoad() {
				delete(props, key)
			}
		}
	}

	// Resolve props values concurrently.
	resolveCtx, cancel := context.WithCancel(r.Context())
	defer cancel()

	type result struct {
		key string
		val any
	}

	resultCh := make(chan result, len(props))
	errCh := make(chan error, 1)

	wg := sync.WaitGroup{}

	for key, val := range props {
		wg.Add(1)

		go func(ctx context.Context, key string, val any) {
			defer wg.Done()

			resolvedVal, err := resolvePropVal(ctx, val)
			if err != nil {
				select {
				case errCh <- fmt.Errorf("resolve prop %q: %w", key, err):
					cancel()
				default:
				}
				return
			}

			select {
			case resultCh <- result{key, resolvedVal}:
			case <-ctx.Done():
			}
		}(resolveCtx, key, val)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for res := range resultCh {
		props[res.key] = res.val
	}

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return props, nil
}

func isPartial(r *http.Request, component string) bool {
	return partialComponentFromRequest(r) == component
}

func getOnlyAndExcept(r *http.Request) (only, except map[string]struct{}) {
	return setOf(onlyFromRequest(r)), setOf(exceptFromRequest(r))
}

func resolvePropVal(ctx context.Context, val any) (_ any, err error) {
	switch proper := val.(type) {
	case Proper:
		val = proper.Prop()
	case TryProper:
		val, err = proper.TryProp()
		if err != nil {
			return nil, err
		}
	case ProperWithContext:
		val = proper.PropWithContext(ctx)
	case TryProperWithContext:
		val, err = proper.TryPropWithContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	switch typed := val.(type) {
	case func() any:
		val = typed()
	case func(ctx context.Context) any:
		val = typed(ctx)
	case func() (any, error):
		val, err = typed()
		if err != nil {
			return nil, fmt.Errorf("closure prop resolving: %w", err)
		}
	case func(ctx context.Context) (any, error):
		val, err = typed(ctx)
		if err != nil {
			return nil, fmt.Errorf("closure prop resolving: %w", err)
		}
	}

	return val, nil
}

func (i *Inertia) resolveEncryptHistory(ctx context.Context) bool {
	encryptHistory, ok := EncryptHistoryFromContext(ctx)
	if ok {
		return encryptHistory
	}
	return i.encryptHistory
}

func (i *Inertia) doInertiaResponse(w http.ResponseWriter, page *page) error {
	pageJSON, err := i.jsonMarshaller.Marshal(page)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	setInertiaInResponse(w)
	setJSONResponse(w)
	setResponseStatus(w, http.StatusOK)

	if _, err = w.Write(pageJSON); err != nil {
		return fmt.Errorf("write bytes to response: %w", err)
	}

	return nil
}

func (i *Inertia) doHTMLResponse(w http.ResponseWriter, r *http.Request, page *page) (err error) {
	// If root template is already created - we'll use it to save some time.
	if i.rootTemplate == nil {
		i.rootTemplate, err = i.buildRootTemplate()
		if err != nil {
			return fmt.Errorf("build root template: %w", err)
		}
	}

	templateData, err := i.buildTemplateData(r, page)
	if err != nil {
		return fmt.Errorf("build template data: %w", err)
	}

	setHTMLResponse(w)

	if err = i.rootTemplate.Execute(w, templateData); err != nil {
		return fmt.Errorf("execute root template: %w", err)
	}

	return nil
}

func (i *Inertia) buildRootTemplate() (*template.Template, error) {
	tmpl := template.New("").Funcs(template.FuncMap(i.sharedTemplateFuncs))
	return tmpl.Parse(i.rootTemplateHTML)
}

func (i *Inertia) buildTemplateData(r *http.Request, page *page) (TemplateData, error) {
	// Defaults.
	inertia, inertiaHead, err := i.buildInertiaHTML(page)
	if err != nil {
		return nil, fmt.Errorf("build inertia html: %w", err)
	}

	templateDataCtx := TemplateDataFromContext(r.Context())

	templateData := make(TemplateData, len(i.sharedTemplateData)+len(templateDataCtx)+2)
	templateData = TemplateData{
		"inertia":     inertia,
		"inertiaHead": inertiaHead,
	}

	// Add shared template data to the result.
	maps.Copy(templateData, i.sharedTemplateData)

	// Add template data from context to the result.
	maps.Copy(templateData, templateDataCtx)

	return templateData, nil
}

func (i *Inertia) buildInertiaHTML(page *page) (inertia, inertiaHead template.HTML, _ error) {
	pageJSON, err := i.jsonMarshaller.Marshal(page)
	if err != nil {
		return "", "", fmt.Errorf("json marshal page into json: %w", err)
	}

	if i.isSSREnabled() {
		inertia, inertiaHead, err = i.htmlContainerSSR(pageJSON)
		if err == nil {
			return inertia, inertiaHead, nil
		}

		i.logger.Printf("ssr rendering error: %s", err)
	}

	return i.htmlContainer(pageJSON)
}

func (i *Inertia) isSSREnabled() bool {
	return i.ssrURL != ""
}

// htmlContainerSSR will send request with json marshaled page payload to ssr render endpoint.
// That endpoint will return head and body html, which will be returned and then rendered.
func (i *Inertia) htmlContainerSSR(pageJSON []byte) (inertia, inertiaHead template.HTML, _ error) {
	url := i.prepareSSRURL()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(pageJSON))
	if err != nil {
		return "", "", fmt.Errorf("new http request: %w", err)
	}
	setJSONRequest(req)

	resp, err := i.ssrHTTPClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("execute http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", "", fmt.Errorf("invalid response status code: %d", resp.StatusCode)
	}

	var ssr struct {
		Head []string `json:"head"`
		Body string   `json:"body"`
	}
	err = i.jsonMarshaller.Decode(resp.Body, &ssr)
	if err != nil {
		return "", "", fmt.Errorf("json decode ssr render response: %w", err)
	}

	inertia = template.HTML(ssr.Body)
	inertiaHead = template.HTML(strings.Join(ssr.Head, "\n"))

	return inertia, inertiaHead, nil
}

func (i *Inertia) prepareSSRURL() string {
	return strings.ReplaceAll(i.ssrURL, "/render", "") + "/render"
}

func (i *Inertia) htmlContainer(pageJSON []byte) (inertia, _ template.HTML, _ error) {
	var sb strings.Builder

	// It doesn't look pretty, but fast!
	sb.WriteString(`<div id="`)
	sb.WriteString(i.containerID)
	sb.WriteString(`" data-page="`)
	template.HTMLEscape(&sb, pageJSON)
	sb.WriteString(`"></div>`)

	return template.HTML(sb.String()), "", nil
}
