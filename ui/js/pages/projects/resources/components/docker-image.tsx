import { FileText } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Docker Image',
    description: 'You can deploy an existing Docker Image from any Registry, without Git.',
    icon: FileText,
    category: 'Docker Based',
};

export default function DockerImage() {
    return (
        <>
            <ResourceCard {...cardProps} />
        </>
    );
}
