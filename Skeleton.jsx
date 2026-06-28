export default function Skeleton({ width = '100%', height = 14, style = {} }) {
  return (
    <div
      className="wr-skel"
      style={{ width, height, ...style }}
    />
  );
}
