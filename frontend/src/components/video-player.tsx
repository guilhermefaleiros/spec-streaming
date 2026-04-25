type Props = {
  manifestUrl: string
}

export function VideoPlayer({ manifestUrl }: Props) {
  return (
    <video controls src={manifestUrl}>
      Your browser does not support the video tag.
    </video>
  )
}
