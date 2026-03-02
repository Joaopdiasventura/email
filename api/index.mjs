import { createTransport } from "nodemailer";

const allowedOrigins = new Set([
  "http://localhost:4200",
  "https://joaopdias-dev.vercel.app",
]);

export default async function handler(req, res) {
  const origin = req.headers.origin;

  if (origin && allowedOrigins.has(origin))
    res.setHeader("Access-Control-Allow-Origin", origin);

  res.setHeader("Vary", "Origin");
  res.setHeader("Access-Control-Allow-Methods", "POST, OPTIONS");
  res.setHeader("Access-Control-Allow-Headers", "Content-Type, Authorization");

  if (req.method == "OPTIONS") return res.status(204).end();

  if (req.method != "POST") {
    res.setHeader("Allow", "POST");
    return res.status(405).json({ error: "METHOD_NOT_ALLOWED" });
  }

  try {
    const body = typeof req.body == "string" ? JSON.parse(req.body) : req.body;
    const subject = String(body?.subject || "").trim();
    const text = body?.text != null ? String(body.text) : undefined;
    const html = body?.html != null ? String(body.html) : undefined;
    const replyTo =
      body?.replyTo != null ? String(body.replyTo).trim() : undefined;

    const to = "joaopdias.dev@gmail.com";

    if (!to || !subject || (!text && !html))
      return res.status(400).json({ error: "INVALID_PAYLOAD" });

    const host = process.env.SMTP_HOST;
    const port = Number(process.env.SMTP_PORT || "587");
    const user = process.env.SMTP_USER;
    const pass = process.env.SMTP_PASS;
    const from = process.env.SMTP_FROM;

    if (!host || !user || !pass || !from)
      return res.status(500).json({ error: "MISSING_SMTP_ENV" });

    const secure =
      String(process.env.SMTP_SECURE || "").toLowerCase() == "true" ||
      port == 465;

    const transporter = createTransport({
      host,
      port,
      secure,
      auth: { user, pass },
    });

    const info = await transporter.sendMail({
      from,
      to,
      subject,
      text,
      html,
      replyTo,
    });

    return res.status(200).json({ ok: true, messageId: info.messageId });
  } catch (e) {
    console.error(e);
    return res.status(500).json({ error: "SEND_FAILED" });
  }
}