import type { FormEvent } from "react";

type RateForm = {
  baseCode: string;
  targetCode: string;
  rate: string;
};

type RatePanelProps = {
  form: RateForm;
  message: string;
  error: string;
  onChange: (nextForm: RateForm) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  normalizeDecimalInput: (value: string) => string;
};

export default function RatePanel({
  form,
  message,
  error,
  onChange,
  onSubmit,
  normalizeDecimalInput,
}: RatePanelProps) {
  return (
    <section className="panel">
      <h2>Rates</h2>
      <form className="stack" onSubmit={onSubmit}>
        <div className="field-row">
          <input
            placeholder="Base"
            value={form.baseCode}
            onChange={(event) =>
              onChange({ ...form, baseCode: event.target.value })
            }
          />
          <input
            placeholder="Target"
            value={form.targetCode}
            onChange={(event) =>
              onChange({ ...form, targetCode: event.target.value })
            }
          />
        </div>
        <input
          placeholder="Rate"
          value={form.rate}
          onChange={(event) =>
            onChange({
              ...form,
              rate: normalizeDecimalInput(event.target.value),
            })
          }
        />
        <button className="primary" type="submit">
          Create rate
        </button>
      </form>
      {message && <p className="success">{message}</p>}
      {error && <p className="error">{error}</p>}
    </section>
  );
}
