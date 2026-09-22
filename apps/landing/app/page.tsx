import { Hero } from "./components/Hero";
import { Thread } from "./components/Thread";
import { Frustrations } from "./components/Frustrations";
import { Weekend } from "./components/Weekend";
import { Capacity } from "./components/Capacity";
import { Method } from "./components/Method";
import { Record } from "./components/Record";
import { Closing } from "./components/Closing";
import { Footer } from "./components/Footer";

export default function Home() {
  return (
    <>
      <main>
        <Hero />
        <div className="relative">
          <Thread />
          <Frustrations />
          <Weekend />
          <Capacity />
          <Method />
          <Record />
        </div>
        <Closing />
      </main>
      <Footer />
    </>
  );
}
