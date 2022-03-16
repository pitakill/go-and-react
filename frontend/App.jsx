const App = () => {
  const [counter, setCounter] = React.useState(0);

  const buttons = [
    { text: "add", fn: () => setCounter(counter + 1) },
    { text: "subtract", fn: () => setCounter(counter - 1) },
  ];

  return (
    <div className="container">
      <h1 style={{ textAlign: "center" }}>{counter}</h1>
      <section className="grid">
        {buttons.map((data, i) => (
          <button key={i} onClick={data.fn} className="secondary">
            {data.text}
          </button>
        ))}
      </section>
    </div>
  );
};

ReactDOM.render(<App />, document.getElementById("root"));
