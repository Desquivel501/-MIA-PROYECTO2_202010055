import { Button } from 'react-bootstrap';
import { Link } from "react-router-dom";

export const Login = () => {

    const ingresar = () => {
        
        

        var user_text = document.getElementById("username").value; 

        if (user_text != ""){
            var pass_text = document.getElementById("password").value; 
            var id_text = document.getElementById("id_particion").value; 
            localStorage.setItem("user", user_text)
        }
      }

    return (
 
        <div className="d-flex fill align-items-center justify-content-center">
            <div class="d-flex flex-column justify-content-center">
                <h1>Iniciar Sesion</h1>
                <form>
                    <div class="mb-3">
                        <label for="exampleInputEmail1" class="form-label">Nombre de Usuario</label>
                        <input type="text" class="form-control" id="username"/>
                    </div>
                    <div class="mb-3">
                        <label for="exampleInputPassword1" class="form-label">Contraseña</label>
                        <input type="text" class="form-control" id="password"/>
                    </div>
                    <div class="mb-3">
                        <label for="exampleInputPassword1" class="form-label">ID de la Particion</label>
                        <input type="text" class="form-control" id="id_particion"/>
                    </div>
                    

                    <div className='row justify-content-end mt-3'>
                        <Button
                            className='col-sm mx-2'
                            onClick={ingresar}
                            as={Link} to="/"
                        >Ingresar</Button>
                    </div>    
                </form>
            </div>
        </div>


    )
  }

  export default {Login}
  