import '#/css/xterm.css';
import { useTerminal } from '@/hooks/use-terminal';
import AppLayout from '@/layouts/app-layout';
import { BreadcrumbItem } from '@/types';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Server',
        href: '/servers',
    },
];

export default function Terminal() {
    const { TerminalDiv } = useTerminal({
        initialContent: 'Hello bro $ ',
    });

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <TerminalDiv className="mx-4 mt-2" />
        </AppLayout>
    );
}
