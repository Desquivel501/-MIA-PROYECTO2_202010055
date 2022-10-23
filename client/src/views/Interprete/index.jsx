
import { Console } from '../../components/Console';
import { Button } from 'react-bootstrap';
import { useState } from 'react';
import { Graphviz } from "graphviz-react";

import './index.css';

export const Interprete = () => {
  const [code, setCode] = useState('');
  const [consoleText, setConsoleText] = useState('');
  const [graph, setGraph] = useState('graph G{}');

  // let graph = "graph G{}"
  let Options = {
    fit: true,
    height: 700,
    width: 1000,
    zoom: true,
  };

  const ejecutar = () => {

    fetch('http://127.0.0.1:5000/consola', {
      method: 'POST',
      body: JSON.stringify({instrucciones:code}),
      headers: {
        'Content-Type':'application/json'
      }
    })
      .then(resp => resp.json())
      .then(data => {

        setGraph(data.reporte)
        setConsoleText(data.resultado)

      })
      .catch(console.error);   
  }

  const clear = () => {
    setConsoleText('');
    setCode('');
  }

  const cargarArchivo = e => {
    e.preventDefault()
    const reader = new FileReader()
    reader.onload = async (e) => { 
      const text = (e.target.result)
      setCode(text)
    };
    reader.readAsText(e.target.files[0])
  };

  return (
    <div>

    <div className="d-flex fill flex-column justify-content-start">

      <div className='row flex-grow-1 mt-3'>
        <Button
              className='col-sm mx-2'
              onClick={() => document.getElementById('fileInput').click()}
        >Cargar Archivo</Button>
        
        <Button
              className='col-sm mx-2'
              onClick={ejecutar}
              variant="success"
            >Ejecutar</Button>
            
        <Button
              className='col-sm mx-2'
              onClick={clear}
              variant="danger"
            >Clear</Button>
      </div>

      <div className='row flex-grow-1 mt-3'>
          <Console code={code} setCode={setCode}></ Console>
          <Console readOnly code={consoleText} setCode={setConsoleText}></Console>
      </div>

      <input id="fileInput" type="file" onChange={cargarArchivo} style={{ display: "none" }} />
    </div>

    <div class='wrapper_graph mt-3' id="graphviz">
        <h4 className="text-light">Reporte</h4>
        <Graphviz 
          classname = "canvas" 
          dot={graph}
          options = {Options}
        />
    </div>

    </div>

    
  )
}
