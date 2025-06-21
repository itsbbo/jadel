import { FileText } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Dockerfile',
    description: 'You can deploy a simple Dockerfile, without Git.',
    icon: FileText,
    category: 'Docker Based',
};

export default function Dockerfile() {
    return (
        <>
            <ResourceCard {...cardProps} />
        </>
    );
}
