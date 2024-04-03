"use server";

export default async function LoadReviews(appId, hours) {
    "use server";
        console.log("loading reviews for app id " + appId + " and hours " + hours);
        try {
        let reviewsReq = await fetch('http://localhost:8000/' + appId + '?hours=' + hours);
        console.log('created fetch');
        let reviews = await reviewsReq.json();
        console.log('got reviews ', reviews);
        return reviews
        } catch (e) {
            throw new Error("Failed to load reviews");
        }
    }