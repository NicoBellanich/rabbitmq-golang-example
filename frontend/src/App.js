import { useState } from "react";

function App() {
  const [msg, setMsg] = useState("");
  const [resp, setResp] = useState("");

  const sendMsg = async () => {
    const res = await fetch(`http://localhost:8080/publish?msg=${msg}`);
    setResp(await res.text());
  };

  return (
    <div className="p-6 text-center">
      <h1>RabbitMQ Demo 🐇</h1>
      <input
        placeholder="Escribí un mensaje"
        value={msg}
        onChange={(e) => setMsg(e.target.value)}
      />
      <button onClick={sendMsg}>Enviar</button>
      <p>{resp}</p>
    </div>
  );
}

export default App;
