import InputError from "@/components/input-error";
import { DialogFooter } from "@/components/shadcn/dialog";
import { Input } from "@/components/shadcn/input";
import { Button } from "@/components/shadcn/button";
import { useForm } from "@inertiajs/react";
import { Dialog, DialogTrigger, DialogContent, DialogTitle, DialogClose } from "@/components/shadcn/dialog";
import { Label } from "@/components/shadcn/label";
import { PlusIcon } from "lucide-react";
import { FormEventHandler } from "react";
import { Textarea } from "@/components/shadcn/textarea";

export default function AddPrivateKey() {
    const { data, setData, post, processing, errors, clearErrors, reset } = useForm({
        name: '',
        description: '',
        publicKey: '',
        privateKey: '',
    });

    const createPrivateKey: FormEventHandler = (e) => {
        e.preventDefault();

        console.log(data);

        post('/private-keys', {
            preserveScroll: true,
            onSuccess: () => closeModal(),
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
                    <DialogTitle>Create Private Key</DialogTitle>
                    <form className="space-y-6" onSubmit={createPrivateKey}>
                        <section className="grid grid-cols-1 gap-4 md:grid-cols-2">
                            <div className="space-y-2">
                                <Label>Name</Label>
                                <Input
                                    id="name"
                                    type="text"
                                    name="name"
                                    value={data.name}
                                    onChange={(e) => setData('name', e.target.value)}
                                    placeholder="Name"
                                />
                                <InputError message={errors.name} />
                            </div>
                            <div className="space-y-2">
                                <Label>Description</Label>
                                <Input
                                    id="description"
                                    type="text"
                                    name="description"
                                    value={data.description}
                                    onChange={(e) => setData('description', e.target.value)}
                                    placeholder="Description"
                                />
                                <InputError message={errors.description} />
                            </div>
                        </section>

                        <section className="space-y-2">
                            <Label className="text-sm font-medium">
                                Private Key
                            </Label>
                            <Textarea
                                id="privateKey"
                                value={data.privateKey}
                                onChange={(e) => setData('privateKey', e.target.value)}
                                placeholder="---BEGIN OPENSSH PRIVATE KEY---"
                            />
                            <InputError message={errors.privateKey} />
                        </section>

                        <section className="space-y-2">
                            <Label className="text-sm font-medium">
                                Public Key
                            </Label>
                            <Textarea
                                id="publicKey"
                                value={data.publicKey}
                                onChange={(e) => setData('publicKey', e.target.value)}
                            />
                            <InputError message={errors.publicKey} />
                        </section>


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