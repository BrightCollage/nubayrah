import { Accordion } from "flowbite-react";
import { NavBar } from "components/NavBar";
import { DefaultFooter } from "components/DefaultFooter";

export default function Home() {
    return (
        <>
            <NavBar></NavBar>
            <div className="grid grid-cols-5">
                <Accordion collapseAll className="col-start-2 col-span-3">
                    <Accordion.Panel>
                        <Accordion.Title>What is Nubayrah?</Accordion.Title>
                        <Accordion.Content>
                            <p className="mb-2 text-gray-500 dark:text-gray-400">
                                Nubayrah is a tool to help you modify and organize your e-books. It helps you configure metadata and extracts information you want as needed.
                            </p>
                            <p className="text-gray-500 dark:text-gray-400">
                                Check out our&nbsp;
                                <a
                                    href="https://github.com/BrightCollage/nubayrah"
                                    className="text-cyan-600 hover:underline dark:text-cyan-500"
                                >
                                    github&nbsp;
                                </a>
                                to learn more!
                            </p>
                        </Accordion.Content>
                    </Accordion.Panel>
                    <Accordion.Panel>
                        <Accordion.Title>Who is working on this Project?</Accordion.Title>
                        <Accordion.Content>
                            <p className="mb-2 text-gray-500 dark:text-gray-400">
                                Two dudes called NotTheBrightestHuman and FartCollage.
                            </p>
                        </Accordion.Content>
                    </Accordion.Panel>
                    <Accordion.Panel>
                        <Accordion.Title>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</Accordion.Title>
                        <Accordion.Content>
                            <p className="mb-2 text-gray-500 dark:text-gray-400">
                                Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed in dignissim urna, vel mattis nulla.
                                Quisque scelerisque dictum ante vel malesuada. Pellentesque euismod dignissim mauris, nec euismod urna ultrices in.
                                Morbi viverra, lectus in cursus vestibulum, nisl diam sollicitudin ante, sed pharetra magna elit vel mauris.
                                Fusce quam mi, pretium in pharetra eget, suscipit id sem. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;
                                Maecenas rhoncus luctus ultricies. Vivamus semper ultrices purus consequat maximus. In malesuada elementum sem ut dictum.
                                Proin suscipit at nisl nec auctor. Proin eget lectus luctus, mattis metus quis, imperdiet urna. Mauris consectetur faucibus orci, ut malesuada metus.
                            </p>
                            <p className="mb-2 text-gray-500 dark:text-gray-400">
                                Etiam id justo odio. Aenean posuere mi finibus mi varius, sed consectetur ante semper. Cras in eleifend elit, ut dictum enim.
                                Phasellus a metus non ipsum malesuada congue a nec nulla. Aliquam magna eros, finibus sit amet venenatis vel, rutrum a elit.
                                Mauris et orci finibus, volutpat quam vel, malesuada neque. Integer auctor sapien nibh, ut scelerisque elit accumsan in.
                                Vivamus porttitor, mauris nec luctus sodales, dui nibh sodales orci, id pharetra nisl quam ac nisi.
                                Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Etiam eget nibh et neque pellentesque iaculis vel in metus.
                                Vestibulum rutrum iaculis mauris, vitae consectetur turpis placerat vel. Integer imperdiet egestas turpis ut elementum.
                                In dui nisl, tincidunt a scelerisque eu, efficitur eu ligula.
                            </p>
                        </Accordion.Content>
                    </Accordion.Panel>
                </Accordion>
            </div>
            <DefaultFooter></DefaultFooter>
        </>
    )
}
