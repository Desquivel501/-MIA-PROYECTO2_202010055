import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';

import { Link } from "react-router-dom";

export const NavbarComponent = () => {
  return (
    <>

      <Navbar bg="dark" variant="dark">
        <Nav className="justify-content-end">
          <Nav.Item className="ms-4 mt-1">
            <Navbar.Brand as={Link} to="/">
              <img
                alt=""
                src="./src/assets/file_logo.png"
                width="30"
                height="30"
                className="d-inline-block align-top"
              />
              {' '}Consola</Navbar.Brand>
            </Nav.Item> 
            <Nav.Item>
            <Nav.Link as={Link} to="/about">Acerca de</Nav.Link>
            </Nav.Item>
        </Nav>
      </Navbar>

    </>
  )
}
