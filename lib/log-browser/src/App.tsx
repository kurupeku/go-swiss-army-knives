import { useState } from 'react'
import Terminal, { ColorMode, TerminalOutput } from 'react-terminal-ui';


function getWindowHeight() {
  const { innerHeight: height } = window;
  return {
    height
  };
}

function App() {
  const { height } = getWindowHeight()
  const termHeight = `${height - 40 - 110}px`
  const [terminalLineData, setTerminalLineData] = useState([
    <TerminalOutput>Welcome to the Log Browser!</TerminalOutput>,
  ]);

  const onInput = (input: string) => {
    if (input === "clear") {
      setTerminalLineData([])
    } else
      setTerminalLineData([...terminalLineData, <TerminalOutput>{input}</TerminalOutput>])
  }

  return (
    <div className="App">
      <main className="uk-margin-top uk-margin-right uk-margin-left uk-margin-bottom">
        <div className="ui-container uk-container-large">
          <Terminal name='Forwarded Logs'
                    colorMode={ColorMode.Dark}
                    onInput={onInput}
                    height={termHeight}
          >
            {terminalLineData}
          </Terminal>
        </div>
      </main>
    </div>
  )
}

export default App
