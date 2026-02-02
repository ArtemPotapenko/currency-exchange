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
  formatRateInput: (value: string) => string;
};

export default function RatePanel({
  form,
  message,
  error,
  onChange,
  onSubmit,
  formatRateInput,
}: RatePanelProps) {
  return (
    <section className="panel">
      <h2>Rates</h2>
      <form className="stack" onSubmit={onSubmit}>
        <div className="field-row">
          <input
            placeholder="Base"
            value={form.baseCode}
            maxLength={3}
            onChange={(event) =>
              onChange({ ...form, baseCode: event.target.value })
            }
          />
          <input
            placeholder="Target"
            value={form.targetCode}
            maxLength={3}
            onChange={(event) =>
              onChange({ ...form, targetCode: event.target.value })
            }
          />
        </div>
        <input
          placeholder="Rate"
          value={form.rate}
          inputMode="decimal"
          onChange={(event) =>
            onChange({
              ...form,
              rate: formatRateInput(event.target.value),
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
