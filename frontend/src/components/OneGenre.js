import { useEffect, useState } from "react";
import { Link, useLocation, useParams } from "react-router-dom"

const OneGenre = () => {
    // We need to get the prop passed to this component
    const location = useLocation();
    const { genreName } = location.state;

    // Set stateful variables
    const [movies, setMovies] = useState([]);

    // Get the ID from url
    let { id } = useParams();

    // useEffec to get list of movies
    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")

        const requestOptions = {
            mehthod: "GET",
            headers: headers,
        }

        fetch(`${process.env.REACT_APP_BACKEND}/movies/genres/${id}`, requestOptions)
            .then(resp => resp.json())
            .then(data => {
                if (data.error) {
                    console.log(data.message)
                } else {
                    setMovies(data);
                }
            })
            .catch(err => {console.log(err)})
    }, [id])

    return (
        <>
            <h2>Genre: {genreName}</h2>
            <hr />

            {movies ? (
            <table className="table table-striped table-hover">
                <thead>
                    <tr>
                        <th>Movie</th>
                        <th>Release Date</th>
                        <th>MPAA Rating</th>
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
                                <td>{m.release_date}</td>
                                <td>{m.mpaa_rating}</td>
                            </tr>
                        ))}
                    </tbody>
                
            </table>
            ): (
                <p>No movies in this genre (yet)!</p>
            )}
        </>
    )
}

export default OneGenre;