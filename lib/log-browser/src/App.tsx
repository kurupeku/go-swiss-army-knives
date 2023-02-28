import { useEffect, useRef, useState } from 'react'
import Terminal, { ColorMode, TerminalOutput } from 'react-terminal-ui';

const url = "ws://localhost:3000/ws"

function getWindowHeight() {
  const { innerHeight: height } = window;
  return {
    height
  };
}

function App() {
  const { height } = getWindowHeight()
  const termHeight = `${height - 40 - 110}px`
  const [lines, setLines] = useState([
    "Welcome to the Log Browser!",
  ]);
  const wsRef = useRef<WebSocket>()
  const linesRef = useRef<string[]>([])

  const onOpen = () => {
    console.log('Connected')
  }

  const onClose = () => {
    console.error('Closed')
  }

  const onMessage = (event: MessageEvent<string>) => {
    const newLines = event.data.split('\n').filter(v => v)
    setLines([...linesRef.current, ...newLines])
    console.log(`Received data: ${event.data}`)
  }

  useEffect(() => {
    wsRef.current = new WebSocket(url)
    wsRef.current.addEventListener("open", onOpen)
    wsRef.current.addEventListener("close", onClose)
    wsRef.current.addEventListener("message", onMessage)

    return () => {
      if (wsRef.current == null || !wsRef.current.CONNECTING) {
        return
      }
      wsRef.current.close()
    }
  }, [])

  useEffect(() => {
    linesRef.current = lines
  }, [lines])

  const onInput = (line: string) => {
    if (line === "clear") {
      setLines([])
    } else
      setLines([...lines, line])
  }

  return (
    <div className="App">
      <main className="uk-margin-top uk-margin-right uk-margin-left uk-margin-bottom">
        <Terminal
          name='Forwarded Logs'
          colorMode={ColorMode.Dark}
          onInput={onInput}
          height={termHeight}
        >
          {lines.map((line, i) => <TerminalOutput key={`${line}-${i}`}>{line}</TerminalOutput>)}
        </Terminal>
      </main>
    </div>
  )
}

export default App
