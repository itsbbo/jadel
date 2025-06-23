import InputError from '@/components/input-error';
import { Button } from '@/components/shadcn/button';
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogTitle, DialogTrigger } from '@/components/shadcn/dialog';
import { Input } from '@/components/shadcn/input';
import { Label } from '@/components/shadcn/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/shadcn/select';
import { PrivateKey } from '@/types/entity';
import { useForm } from '@inertiajs/react';
import { PlusIcon } from 'lucide-react';
import { FormEventHandler, useEffect, useState } from 'react';

export default function AddServer() {
    const [privateKeys, setPrivateKeys] = useState<PrivateKey[]>([]);

    const { data, setData, post, processing, errors, setError, clearErrors, reset } = useForm({
        name: '',
        description: '',
        ip: '127.0.0.1',
        port: 22,
        user: 'root',
        privateKeyId: '',
    });

    useEffect(() => {
        fetch('/private-keys/all/json')
            .then((res) => res.json())
            .then((data) => setPrivateKeys(data))
            .catch((_) => setError('privateKeyId', 'Cannot get your private keys. Try again later.'));
    }, []);

    const createServer: FormEventHandler = (e) => {
        e.preventDefault();

        console.log(data);

        post('/servers', {
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
                    <DialogTitle>Create Server</DialogTitle>
                    <form className="space-y-6" onSubmit={createServer}>
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

                        <section className="grid grid-cols-1 gap-4 md:grid-cols-2">
                            <div className="space-y-2">
                                <Label>IP Address</Label>
                                <Input
                                    id="ip"
                                    type="text"
                                    name="ip"
                                    value={data.ip}
                                    onChange={(e) => setData('ip', e.target.value)}
                                    placeholder="Name"
                                />
                                <InputError message={errors.ip} />
                            </div>
                            <div className="space-y-2">
                                <Label>Port</Label>
                                <Input
                                    id="port"
                                    type="number"
                                    name="port"
                                    value={data.port}
                                    onChange={(e) => setData('port', Number(e.target.value))}
                                    placeholder="22"
                                />
                                <InputError message={errors.port} />
                            </div>
                        </section>

                        <section className="space-y-2">
                            <Label htmlFor="user" className="text-sm font-medium">
                                User
                            </Label>
                            <Input id="user" value={data.user} className="border-gray-600" onChange={(e) => setData('user', e.target.value)} />
                            <InputError message={errors.user} />
                        </section>

                        <section className="space-y-2">
                            <Label className="text-sm font-medium">Private Key</Label>
                            <Select value={data.privateKeyId} onValueChange={(e) => setData('privateKeyId', e)}>
                                <SelectTrigger className="w-full">
                                    <SelectValue placeholder="Select private key" />
                                </SelectTrigger>
                                <SelectContent>
                                    {privateKeys.map((key) => (
                                        <SelectItem key={key.id} value={key.id}>
                                            {key.name}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                            <InputError message={errors.privateKeyId} />
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
