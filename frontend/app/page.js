import Image from "next/image";

export default function Home() {
  return (
    <main>
      <input type="text" name="appId" />
      <button type="button" onClick={loadReviews}>Load Reviews</button>
    </main>
  );
}
