'use client';
import Review from "./review";

export default function ReviewList({ reviews }) {
  const emptyList = reviews.length < 1
  return (
    <div className="grid grid-cols-12 gap-2 auto">
      <div className="col-start-3 col-span-8">
        { emptyList && <h2>No new reviews</h2> }
        { reviews.map((review) => (
            <Review key={review.id} review={review} />
        ))}
      </div>
    </div>
  );
}