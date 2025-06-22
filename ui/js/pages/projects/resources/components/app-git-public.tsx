import InputError from '@/components/input-error';
import { Button } from '@/components/shadcn/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/shadcn/dialog';
import { Input } from '@/components/shadcn/input';
import { Label } from '@/components/shadcn/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/shadcn/select';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/shadcn/tooltip';
import { Environment } from '@/types/entity';
import { useForm } from '@inertiajs/react';
import { GitBranch, Info } from 'lucide-react';
import { FormEventHandler, useEffect, useState } from 'react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Public Repository',
    description: 'You can deploy any kind of public repositories from the supported git providers.',
    icon: GitBranch,
    category: 'Git Based',
    isPopular: true,
};

type FormData = {
    repositoryURL: string;
    branch: string;
    buildTool: string;
    baseDir: string;
    target: string;
    port: number;
    isStaticSite: boolean;
};

interface Props {
    env: Environment;
}

export default function GitPublic({ env }: Props) {
    const [open, setOpen] = useState(false);
    const { data, setData, post, reset, errors, processing } = useForm<FormData>({
        repositoryURL: '',
        branch: 'main',
        buildTool: 'Nixpacks',
        baseDir: '/',
        target: '/',
        port: 3000,
        isStaticSite: false,
    });
    const submit: FormEventHandler = (e) => {
        e.preventDefault();

        post(`/projects/${env.project.id}/environments/${env.id}`, {
            preserveScroll: true,
            onFinish: () => reset(),
        });
    };

    const [targetName, setTargetName] = useState('Build Directory');
    const [targetTooltipContent, setTargetTooltipContent] = useState(
        'The directory where your application is built. Ignore this if you dont have static/asset/binary executable as an output.',
    );

    useEffect(() => {
        switch (data.buildTool) {
            case 'Dockerfile':
                setTargetName('Dockerfile Location');
                setTargetTooltipContent('The location of your Dockerfile. Default is the root of your repository.');
                setData('target', '/Dockerfile');
                break;
            case 'Docker Compose':
                setTargetName('Docker Compose Location');
                setTargetTooltipContent('The location of your Docker Compose file. Default is the root of your repository.');
                setData('target', '/docker-compose.yml');
                break;
            default:
                setTargetName('Build Directory');
                if (data.target !== '/') {
                    setData('target', '/');
                }
                break;
        }
    }, [data.buildTool]);

    return (
        <>
            <Dialog open={open} onOpenChange={setOpen}>
                <DialogContent className="sm:max-w-lg">
                    <DialogHeader>
                        <DialogTitle>Create a New Git Based Resource</DialogTitle>
                        <DialogDescription>Deploy any public Git repositories.</DialogDescription>
                    </DialogHeader>

                    <form className="space-y-6" onSubmit={submit}>
                        <div className="space-y-2">
                            <Label htmlFor="repository-url" className="flex items-center gap-2">
                                Repository URL (https://)
                                <span className="text-orange-500">*</span>
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Info className="text-muted-foreground h-4 w-4 cursor-help" />
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>Enter the URL of your Git repository (e.g., GitHub, GitLab)</p>
                                    </TooltipContent>
                                </Tooltip>
                            </Label>
                            <div className="flex gap-2">
                                <Input
                                    type="url"
                                    placeholder="https://github.com/user/repo"
                                    value={data.repositoryURL}
                                    onChange={(e) => setData('repositoryURL', e.target.value)}
                                    className="flex-1"
                                />
                                <InputError className="mt-2" message={errors.repositoryURL} />
                            </div>
                        </div>

                        <div className="w-full space-y-2">
                            <Label htmlFor="build-pack" className="flex items-center gap-2">
                                Buildtool
                                <span className="text-orange-500">*</span>
                            </Label>
                            <Select value={data.buildTool} onValueChange={(e) => setData('buildTool', e)}>
                                <SelectTrigger className="w-full">
                                    <SelectValue placeholder="Select build pack" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="Nixpacks">Nixpacks</SelectItem>
                                    <SelectItem value="Buildpacks">Buildpacks</SelectItem>
                                    <SelectItem value="Dockerfile">Dockerfile</SelectItem>
                                    <SelectItem value="Docker Compose">Docker Compose</SelectItem>
                                </SelectContent>
                            </Select>
                            <InputError className="mt-2" message={errors.buildTool} />
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div className="w-full space-y-2">
                                <Label htmlFor="branch" className="flex items-center gap-2">
                                    Branch
                                    <Tooltip>
                                        <TooltipTrigger asChild>
                                            <Info className="text-muted-foreground h-4 w-4 cursor-help" />
                                        </TooltipTrigger>
                                        <TooltipContent>
                                            <p>Branch to use for deployment</p>
                                        </TooltipContent>
                                    </Tooltip>
                                </Label>
                                <Input id="branch" value={data.branch} onChange={(e) => setData('branch', e.target.value)} />
                                <InputError className="mt-2" message={errors.branch} />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="port" className="flex items-center gap-2">
                                    Port
                                    <Tooltip>
                                        <TooltipTrigger asChild>
                                            <Info className="text-muted-foreground h-4 w-4 cursor-help" />
                                        </TooltipTrigger>
                                        <TooltipContent>
                                            <p>The port application will be exposed on</p>
                                        </TooltipContent>
                                    </Tooltip>
                                </Label>
                                <Input id="port" type="number" value={data.port} onChange={(e) => setData('port', Number(e.target.value))} />
                                <InputError className="mt-2" message={errors.port} />
                            </div>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="base-directory" className="flex items-center gap-2">
                                Base Directory
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Info className="text-muted-foreground h-4 w-4 cursor-help" />
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>The base directory/root of your application</p>
                                    </TooltipContent>
                                </Tooltip>
                            </Label>
                            <Input id="base-directory" value={data.baseDir} onChange={(e) => setData('baseDir', e.target.value)} />
                            <InputError className="mt-2" message={errors.baseDir} />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="base-directory" className="flex items-center gap-2">
                                {targetName}
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Info className="text-muted-foreground h-4 w-4 cursor-help" />
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>{targetTooltipContent}</p>
                                    </TooltipContent>
                                </Tooltip>
                            </Label>
                            <Input id="base-directory" value={data.target} onChange={(e) => setData('target', e.target.value)} />
                            <InputError className="mt-2" message={errors.target} />
                        </div>

                        <div className="flex items-center space-x-2">
                            <input
                                type="checkbox"
                                id="static-site"
                                checked={data.isStaticSite}
                                onChange={(e) => setData('isStaticSite', e.target.checked)}
                                className="h-4 w-4 cursor-pointer"
                            />
                            <Label htmlFor="static-site" className="flex items-center gap-2">
                                Is it a static site?
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Info className="text-muted-foreground h-4 w-4 cursor-help" />
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>
                                            If your application is a static site or the final build assets should be served as a static site, enable
                                            this.
                                        </p>
                                    </TooltipContent>
                                </Tooltip>
                            </Label>
                        </div>

                        <div className="flex justify-end gap-3">
                            <Button onClick={() => setOpen(false)} disabled={processing}>
                                Create
                            </Button>
                        </div>
                    </form>
                </DialogContent>
            </Dialog>

            <ResourceCard {...cardProps} onClick={() => setOpen(true)} />
        </>
    );
}
