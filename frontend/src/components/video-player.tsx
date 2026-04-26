type Props = {
  manifestUrl: string
}

export function VideoPlayer({ manifestUrl }: Props) {
  return (
    <div style={{ marginTop: '20px' }}>
      <video 
        controls 
        src={manifestUrl}
        style={{ 
          width: '100%', 
          maxWidth: '800px',
          borderRadius: '8px'
        }}
      >
        Your browser does not support the video tag.
      </video>
    </div>
  )
}
