'use client';

import AppInput from "./app-input";
import ReviewList from "./review-list";
import LoadReviews from "./load-reviews";
import { useState } from "react";

export default function AppReviews() {
    // let reviews = [];
    const [reviews, setReviews] = useState([]);
    function loadReviews(appId) {
        LoadReviews(appId).then((reviews) => {
                setReviews(reviews);
        });
    }

    return (
        <div>
            <AppInput loadReviews={loadReviews} />
            <ReviewList reviews={reviews} />
        </div>
    );
}