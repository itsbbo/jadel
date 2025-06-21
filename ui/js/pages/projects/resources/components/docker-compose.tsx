import { FileText } from 'lucide-react';
import { ResourceCard } from '../../components/resource-card';

const cardProps = {
    title: 'Docker Compose Empty',
    description: 'You can deploy complex application easily with Docker Compose, without Git.',
    icon: FileText,
    category: 'Docker Based',
};

export default function DockerCompose() {
    return <ResourceCard {...cardProps} />;
}
