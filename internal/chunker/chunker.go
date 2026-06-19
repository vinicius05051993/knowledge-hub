package chunker

func Split(text string) []string {
    const size = 1000

    var chunks []string

    for len(text) > size {
        chunks = append(chunks, text[:size])
        text = text[size:]
    }

    if text != "" {
        chunks = append(chunks, text)
    }

    return chunks
}