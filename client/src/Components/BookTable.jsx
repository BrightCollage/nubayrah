import { Modal, Table, Dropdown, Alert } from "flowbite-react";
import { HiOutlineTrash, HiDownload, HiInformationCircle } from "react-icons/hi";
import { useState, useEffect } from "react";
import axios from "axios";

export function BookTable() {
    // Book data
    const [books, setBooks] = useState([])

    // State to check if reload is necessary.
    const [reload, setReload] = useState(true)

    // Notification states.
    const [isNotification, setIsNotification] = useState(false)
    const [response, setResponse] = useState([])

    // GET request to /books to retrieve all book data.
    // Refreshes after a delete.
    useEffect(() => {
        if (reload) {
            const fetchData = async () => {
                await axios.get('http://localhost:5050/books/')
                    .then(res => {
                        console.log(res);
                        setBooks(res.data)
                    })
                    .catch(err => {
                        console.log(err.response)
                    })
            }
            fetchData()
        }
        setReload(false)

    }, [reload]) // Reloads on state change.

    return (
        <div className="">
            <Modal show={isNotification} >
                <Alert color={response["status"] === 204 ? "info" : "failure"} icon={HiInformationCircle} onDismiss={() => { setIsNotification(false); setReload(true) }}>
                    <span className="font-medium">Info: </span> Delete successful!
                </Alert>
            </Modal>
            <Table striped>
                <Table.Head>
                    <Table.HeadCell>Title</Table.HeadCell>
                    <Table.HeadCell>Author</Table.HeadCell>
                    <Table.HeadCell>Publish Date</Table.HeadCell>
                    <Table.HeadCell>id</Table.HeadCell>
                    <Table.HeadCell>
                        <span className="sr-only">Edit</span>
                    </Table.HeadCell>
                </Table.Head>
                <Table.Body className="divide-y">
                    {/* Function to grab all books from GET request and create table */}
                    {books.map(book => (
                        <Table.Row className="bg-white dark:border-gray-700 dark:bg-gray-800" key={book['id']}>
                            <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                                {book['title']}
                            </Table.Cell>
                            <Table.Cell>{book['author']}</Table.Cell>
                            <Table.Cell>{book['pubDate']}</Table.Cell>
                            <Table.Cell>{book['id']}</Table.Cell>
                            <Table.Cell>
                                <BookRowDropDown id={book['id']}
                                    setIsNotification={setIsNotification}
                                    setResponse={setResponse}></BookRowDropDown>
                            </Table.Cell>
                        </Table.Row>
                    ))}
                </Table.Body>
            </Table>
        </div>
    );
}

// Dropdown component for book row.
const BookRowDropDown = (props) => {
    const { id, setIsNotification, setResponse } = props

    const openInNewTab = (url) => {
        const newWindow = window.open(url, '_blank', 'noopener,noreferrer')
        if (newWindow) newWindow.opener = null
    }

    const deleteBookRequest = async (e) => {
        setIsNotification(true)
        await axios.delete('http://localhost:5050/books/' + id)
            .then(res => {
                console.log(res)
                setResponse(res)
            })
            .catch(err => {
                console.log(err)
            })

    }

    return (
        <>
            <Dropdown label="" inline>
                <Dropdown.Item icon={HiDownload} onClick={() => openInNewTab('http://localhost:5050/books/' + id + '/cover')}>Cover</Dropdown.Item>
                <Dropdown.Divider />
                <Dropdown.Item icon={HiOutlineTrash} onClick={deleteBookRequest}>Delete</Dropdown.Item>
            </Dropdown>
        </>
    );
}
