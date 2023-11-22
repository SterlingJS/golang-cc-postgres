import logo from './logo.svg';
import './App.css';
import { useState, useEffect, useRef } from 'react';

function App() {
  const [message, setMessage] = useState('')
  const [replicas, setReplicas] = useState(1)
  const [sent, setSent] = useState(0)
  const [failed, setFailed] = useState(0)
  const stateRef = useRef();
  const MAX_MESSAGES = 9999;
  stateRef.success = sent;
  stateRef.fail = failed;

  var producerBackendURI = process.env.REACT_APP_PRODUCER_BACKEND
  if (producerBackendURI != null && producerBackendURI !== "") {
    console.log("Producer backend is set to:")
    console.log(producerBackendURI)
  } else {
    console.log("Producer Backend is not set, using hard coded value")
    producerBackendURI = "https://keda-producer.fmz-c-x-app-aks-01.corp.fmglobal.com/enqueue"
    console.log(process.env.REACT_APP_REACT_APP_PRODUCER_BACKEND)
  }

  const onFormSubmit = async (e) => {
    e.preventDefault();

    for (let i = 0; i < replicas; i++) {
      sendMessage()
    }
  }

  const safelySetReplicas = (val) => {
    if (val > MAX_MESSAGES) {
      setReplicas(MAX_MESSAGES)
    } else {
      setReplicas(val)
    }
  }

  const sendMessage = async () => {
    try {
      const data = { message: message }
      await fetch(producerBackendURI, {
        method: "POST", // *GET, POST, PUT, DELETE, etc.
        mode: "no-cors", // no-cors, *cors, same-origin
        cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
        credentials: "same-origin", // include, *same-origin, omit
        headers: {
          "Content-Type": "application/json",
        },
        redirect: "follow", // manual, *follow, error
        referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: JSON.stringify(data), // body data type must match "Content-Type" header
      })
      setSent((sent) => sent + 1)

    } catch (err) {
      console.error(err)
      setFailed((failed) => failed + 1)
    } 
  }
  
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Enter a message to see KEDA scale!
        </p>

        <form className="App-input-container" onSubmit={onFormSubmit}>
          <div className="App-label">
            <div className="App-label-text">Message:  </div>
            <input className="App-input" type="text" value={message} onChange={e => setMessage(e.currentTarget.value)} required />
          </div>
          <div className="App-label">
            <div className="App-label-text">Replicas: </div>
            <input className="App-input" type="number" value={replicas} onChange={e => safelySetReplicas(e.currentTarget.value)} required />
          </div>
          <button type="submit" >Send</button>
        </form>
        <div>Successful Messages: {stateRef.success}</div>
        <div>Failed Messages: {stateRef.fail}</div>
      </header>
    </div>
  );
}

export default App;
