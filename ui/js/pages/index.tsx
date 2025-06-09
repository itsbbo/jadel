import { Head } from "@inertiajs/react"

interface Props {
    some: string
}

export default function Index({ some }: Props) {
    return (
        <div>
            <Head title="Home" />
            <p>Hello World {some}</p>
        </div>
    )
}