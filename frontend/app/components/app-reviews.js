'use client';

import AppInput from "./app-input";
import ReviewList from "./review-list";
import LoadReviews from "./load-reviews";
import { useState } from "react";

export default function AppReviews() {
    // let reviews = [];
    const [reviews, setReviews] = useState([]);
    const [hasPressedLoad, setHasPressedLoad] = useState(false);
    const [error, setError] = useState("")
    function loadReviews(appId) {
        LoadReviews(appId).then((reviews) => {
            setError("");
            setHasPressedLoad(true);
            setReviews(reviews);
        }, (error) => {
            console.log("setting error", error)
            setError(error.message);
            setHasPressedLoad(false);
        });
    }

    return (
        <div>
            <AppInput loadReviews={loadReviews} />
            { error.length > 0 && 
            <div className="grid grid-cols-12 gap-2 auto my-8">
                <div className="col-start-3 col-span-8 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                    <strong className="font-bold">Error!</strong>
                    <span className="block">{error}</span>
                </div> 
            </div>
            }
            { hasPressedLoad && <ReviewList reviews={reviews} /> }
        </div>
    );
}