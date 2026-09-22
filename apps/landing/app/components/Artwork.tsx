/** One piece of the still-life family: bone plaster, navy ground, window light. */
export function Artwork({ src, alt, className = "" }: { src: string; alt: string; className?: string }) {
  return (
    <figure className={`overflow-hidden rounded-[6px] bg-[#0e1626] ${className}`}>
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={src}
        alt={alt}
        width={1024}
        height={1280}
        loading="lazy"
        decoding="async"
        className="aspect-[5/4] h-full w-full object-cover object-[50%_84%] lg:aspect-[4/5] lg:object-center"
      />
    </figure>
  );
}
