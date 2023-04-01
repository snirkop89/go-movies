import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Input from "./form/Input";

const GraphQL = () => {
    const [movies, setMovies] = useState([]);
    const [searchTerm, setSearchTerm] = useState("");
    const [fullList, setFullList] = useState([]);

    // Perform a search
    const performSearch = () => {
        const payload = `
        {
            search(titleContains: "${searchTerm}") {
                id
                title
                runtime
                release_date
                mpaa_rating
            }
        }`;

        const headers = new Headers();
        headers.append("Content-Type", "application/graphql");

        const requestOptions = {
            method: "POST",
            body: payload,
            headers: headers,
        }

        fetch(`${process.env.REACT_APP_BACKEND}/graph`, requestOptions)
            .then(resp => resp.json())
            .then(data => {
                let theList = Object.values(data.data.search);
                setMovies(theList);
            })
            .catch(err => {console.log(err)})
    }

    const handleChange = (event) => {
        event.preventDefault();

        let value = event.target.value;
        setSearchTerm(value);

        if (value.length > 2) {
            performSearch();
        } else {
            setMovies(fullList);
        }
    }

    // useEffect
    useEffect(() => {
        const payload = `
        {
            list {
                id
                title
                runtime
                release_date
                mpaa_rating
            }
        }`;

        const headers = new Headers();
        headers.append("Content-Type", "application/graphql");
        
        const requestOptions = {
            method: "POST",
            headers: headers,
            body: payload,
        }

        fetch(`${process.env.REACT_APP_BACKEND}/graph`, requestOptions)
            .then(resp => resp.json())
            .then(data => {
                let theList = Object.values(data.data.list);
                setMovies(theList);
                setFullList(theList);
            })
            .catch(err => {console.log(err)})
    }, [])

    return (
        <div className="">
            <h2>GraphQL</h2>
            <hr />

            <form onSubmit={handleChange}>
                <Input
                    title={"Search"}
                    type={"search"}
                    name={"search"}
                    className={"form-control"}
                    value={searchTerm}
                    onChange={handleChange}
                />
            </form>

            {movies ? (
                <table className="table table-striped table-hover">
                    <thead>
                        <tr>
                            <th>Movie</th>
                            <th>Release Date</th>
                            <th>Rating</th>
                        </tr>
                    </thead>
                    <tbody>
                        {movies.map(m => (
                            <tr key={m.id}>
                                <td>
                                    <Link to={`/movies/${m.id}`}>
                                        {m.title}
                                    </Link>
                                </td>
                                <td>{new Date(m.release_date).toLocaleDateString()}</td>
                                <td>{m.mpaa_rating}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            ) : (
                <p>No movies (yet)!</p>
            )}
        </div>

    )
}

export default GraphQL;