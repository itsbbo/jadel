import InputError from '@/components/input-error';
import { Button } from '@/components/shadcn/button';
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogTitle, DialogTrigger } from '@/components/shadcn/dialog';
import { Input } from '@/components/shadcn/input';
import { Label } from '@/components/shadcn/label';
import { useForm } from '@inertiajs/react';
import { PlusIcon } from 'lucide-react';
import { FormEventHandler, useRef } from 'react';

export default function AddProjects() {
    const nameInput = useRef<HTMLInputElement>(null);
    const descriptionInput = useRef<HTMLInputElement>(null);

    const { data, setData, post, processing, errors, clearErrors, reset } = useForm({
        name: '',
        description: '',
    });

    const createProject: FormEventHandler = (e) => {
        e.preventDefault();

        post('/projects', {
            preserveScroll: true,
            onSuccess: () => closeModal(),
            onError: () => nameInput.current?.focus(),
            onFinish: () => reset(),
        });
    };

    const closeModal = () => {
        clearErrors();
        reset();
    };

    return (
        <>
            <Dialog>
                <DialogTrigger asChild>
                    <Button>
                        <PlusIcon /> Add
                    </Button>
                </DialogTrigger>
                <DialogContent>
                    <DialogTitle>Create Project</DialogTitle>
                    <DialogDescription>
                        New project will have a default <b>production</b> environment.
                    </DialogDescription>
                    <form className="space-y-6" onSubmit={createProject}>
                        <div className="grid gap-2">
                            <Label className="sr-only">Name</Label>
                            <Input
                                id="name"
                                type="text"
                                name="name"
                                ref={nameInput}
                                value={data.name}
                                onChange={(e) => setData('name', e.target.value)}
                                placeholder="Name"
                            />

                            <InputError message={errors.name} />
                        </div>

                        <div className="grid gap-2">
                            <Label className="sr-only">Description</Label>

                            <Input
                                id="description"
                                type="text"
                                name="description"
                                ref={descriptionInput}
                                value={data.description}
                                onChange={(e) => setData('description', e.target.value)}
                                placeholder="Description"
                            />

                            <InputError message={errors.description} />
                        </div>

                        <DialogFooter className="gap-2">
                            <DialogClose asChild>
                                <Button variant="secondary" onClick={closeModal}>
                                    Cancel
                                </Button>
                            </DialogClose>

                            <Button disabled={processing} asChild>
                                <button type="submit">Create</button>
                            </Button>
                        </DialogFooter>
                    </form>
                </DialogContent>
            </Dialog>
        </>
    );
}
