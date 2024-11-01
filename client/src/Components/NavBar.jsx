
import { Link, useLocation } from "react-router-dom";
import { DarkThemeToggle, Navbar } from "flowbite-react";
import nubayrahIcon from '../Assets/nubayrah.svg'

export function NavBar() {

    const location = useLocation(); // Once ready, returns 'window.location'

    return (
        <Navbar fluid rounded>
            <Navbar.Brand href="https://github.com/BrightCollage/nubayrah">
                <img src={nubayrahIcon} className="mr-3 h-6 sm:h-9" alt="Flowbite React Logo" />
                <span className="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Nubayrah</span>
            </Navbar.Brand>
            <div className="flex md:order-2">
                <DarkThemeToggle className="flex self-end" />
                <Navbar.Toggle />
            </div>
            <Navbar.Collapse>
                <Navbar.Link href="/" active={location.pathname === "/"}>Home</Navbar.Link>
                <Navbar.Link href="/library" active={location.pathname === "/library"}>Library</Navbar.Link>
                <Navbar.Link href="/upload" active={location.pathname === "/upload"}>Upload</Navbar.Link>
                {/* <Navbar.Link href="#" active={url === "/" ?" active" : ""}>About</Navbar.Link> */}
            </Navbar.Collapse>
        </Navbar>
    );
}