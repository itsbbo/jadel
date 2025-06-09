interface Props {
    some: string
}

export default function Bar({ some }: Props) {
    return (
        <div>
            <p>Hello Some {some}</p>
        </div>
    )
}