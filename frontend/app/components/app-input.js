import ReviewList from "@/components/review-list";

export default function AppInput({ params }) {
  async function loadReviews(event) {
    event.preventDefault();
    console.log(event);

    const response = await fetch('localhost:8000/' + params.appId);
    const reviewData = await response.json();
  }

  return (
    <div>
      <input type="text" name="appId" />
      <button type="button" onClick={loadReviews}>Load Reviews</button>
      <ReviewList reviews={reviewData} />
    </div>
  )
}