
import { Dropdown } from "flowbite-react";
import { HiOutlineTrash, HiDownload } from "react-icons/hi";

export function RowDropDown(props) {
    const { id } = props

    const openInNewTab = (url) => {
        const newWindow = window.open(url, '_blank', 'noopener,noreferrer')
        if (newWindow) newWindow.opener = null
    }

    return (
        <Dropdown label="" inline>
            <Dropdown.Item icon={HiDownload} onClick={() => openInNewTab('http://localhost:5050/books/' + id + '/cover')}>Cover</Dropdown.Item>
            <Dropdown.Divider />
            <Dropdown.Item icon={HiOutlineTrash}>Delete</Dropdown.Item>
        </Dropdown>
    );
}
