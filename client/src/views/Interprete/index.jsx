
import { Console } from '../../components/Console';
import { Button } from 'react-bootstrap';
import { useState } from 'react';
import { Graphviz } from "graphviz-react";
import './index.css';
import swal from 'sweetalert';
import { Link } from "react-router-dom";


export const Interprete = () => {
  const [code, setCode] = useState('');
  const [consoleText, setConsoleText] = useState('');
  const [graph, setGraph] = useState('graph G{}');
  const [state, setState] = useState(0);
  

  const user = localStorage.getItem("user")

  // let graph = "graph G{}"
  let Options = {
    fit: true,
    height: 700,
    width: 1000,
    zoom: true,
  };

  const logout = () => {
    localStorage.setItem("user", "")
    // if (window.confirm("¿Desea cerrar sesión?")) {
    //   localStorage.setItem("user", "")
    //   setState(state+1)
    // }
    swal({
      title: "¿Desea cerrar sesion?",
      icon: "warning",
      buttons: true,
      dangerMode: true,
    })
    .then((willLogout) => {
      if (willLogout) {

        fetch('http://18.118.206.82:5000/logout', {
          method: 'POST',
          body: JSON.stringify({"mensaje":""}),
          headers: {
              'Content-Type':'application/json'
          }
          })
          .then(resp => resp.json())
          .then(data => {
              
              if(data.mensaje != ""){
                    swal(data.mensaje, {
                      icon: "error",
                    });

              } else {

                swal("Se ha cerrrado la sesion", {
                  icon: "success",
                });
                localStorage.setItem("user", "")
                setState(state+1)
              }

              

          })
          .catch(console.error); 


        
      } 
    });
    
  }

  const ejecutar = () => {

    fetch('http://18.118.206.82:5000/consola', {
      method: 'POST',
      body: JSON.stringify({instrucciones:code}),
      headers: {
        'Content-Type':'application/json'
      }
    })
      .then(resp => resp.json())
      .then(data => {

        console.log(data.reporte)

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

  
  
  // alert(user)

  return (
    <div>

    <div className="d-flex fill flex-column justify-content-start">

    {user != "" &&        
      <div className='row justify-content-end mt-3'>
      <Button
              style={{width:"15%"}}
              className='mx-2'
              variant="warning"
              onClick={logout}
        >Cerrar Sesion</Button>
      </div>    
    }

  {user == "" &&        
      <div className='row justify-content-end mt-3'>
      <Button
              style={{width:"15%"}}
              className='mx-2'
              variant="warning"
              as={Link} to="/login"
        >Iniciar Sesion</Button>
      </div>    
    }

    
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


