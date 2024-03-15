'use client';
import StarItem from "./star";

export default function Review({ review }) {
    const stars = [...Array(5)].map(x => 0);
    const date = new Date(review.updated);
    return (
        <div className="max-w-sm w-full lg:max-w-full my-8">
            <div className="border border-gray-400 lg:border-gray-400 bg-white rounded p-4 justify-between leading-normal">
                <div className="mb-8 flex items-center">
                    <StarItem filled={1 <= review.rating} />
                    <StarItem filled={2 <= review.rating} />
                    <StarItem filled={3 <= review.rating} />
                    <StarItem filled={4 <= review.rating} />
                    <StarItem filled={5 <= review.rating} />
                </div>
                <div className="mb-8">
                    <p className="text-gray-700 text-base">{review.content}</p>
                </div>
                <div className="flex items-center">
                    <div className="text-sm">
                        <p className="text-gray-900 leading-none">{review.author.name}</p>
                        <p className="text-gray-600">{date.toLocaleString()}</p>
                    </div>
                </div>
            </div>
        </div>
    );
}