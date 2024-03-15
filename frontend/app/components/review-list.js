'use client';
import Review from "./review";

export default function ReviewList({ reviews }) {

  return (
    <div className="grid grid-cols-12 gap-2 auto">
      <div className="col-start-3 col-span-8">
        { reviews.map((review) => (
            <Review key={review.id} review={review} />
        ))}
      </div>
    </div>
  );
}