export const withPreventDefault = (callback: () => void) => {
  return (event: React.MouseEvent) => {
    event.preventDefault();
    callback();
  };
};
