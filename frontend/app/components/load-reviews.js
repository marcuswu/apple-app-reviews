"use server";

export default async function LoadReviews(appId) {
    "use server";
        console.log("loading reviews");
        try {
        let reviewsReq = await fetch('http://localhost:8000/' + appId);
        console.log('created fetch');
        let reviews = await reviewsReq.json();
        console.log('got reviews ', reviews);
        return reviews
        } catch (e) {
            throw new Error("Failed to load reviews");
        }
    }