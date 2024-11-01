import { Button, Table } from "flowbite-react";
import { useState, useEffect } from "react";
import axios from "axios";
import { RowDropDown } from "./RowDropDown";

export function BookTable() {

    const [books, setBooks] = useState([])

    // GET request to /books to retrieve all book data.
    useEffect(() => {
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
    }, [])


    return (
        <div className="overflow-x-auto">
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
                                <RowDropDown id={book['id']}></RowDropDown>
                            </Table.Cell>
                        </Table.Row>
                    ))}
                </Table.Body>
            </Table>
        </div>
    );
}
