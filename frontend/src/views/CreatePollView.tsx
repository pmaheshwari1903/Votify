import React, { useState } from 'react';
import { Button } from '../components/ui/Button';
import { apiClient } from '../lib/api';
import type { Poll } from '../types';

interface CreatePollViewProps {
  onSuccess: (poll: Poll) => void;
  onCancel: () => void;
}

export const CreatePollView: React.FC<CreatePollViewProps> = ({ onSuccess, onCancel }) => {
  const [question, setQuestion] = useState('');
  const [description, setDescription] = useState('');
  const [options, setOptions] = useState<string[]>(['Option 1', 'Option 2']);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleOptionChange = (index: number, value: string) => {
    const updated = [...options];
    updated[index] = value;
    setOptions(updated);
  };

  const handleAddOption = () => {
    if (options.length >= 10) return;
    setOptions([...options, `Option ${options.length + 1}`]);
  };

  const handleRemoveOption = (index: number) => {
    if (options.length <= 2) return;
    const updated = options.filter((_, i) => i !== index);
    setOptions(updated);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const trimmedQuestion = question.trim();
    if (trimmedQuestion.length < 5) {
      setError('Question must be at least 5 characters long.');
      return;
    }

    const trimmedOptions = options.map((opt) => opt.trim()).filter((opt) => opt.length > 0);
    if (trimmedOptions.length < 2) {
      setError('Please provide at least 2 non-empty options.');
      return;
    }

    // Check duplicate options client-side for UX
    const lowerSet = new Set(trimmedOptions.map((o) => o.toLowerCase()));
    if (lowerSet.size !== trimmedOptions.length) {
      setError('Duplicate options are not allowed.');
      return;
    }

    setLoading(true);
    try {
      const poll = await apiClient.post<Poll>('/polls', {
        question: trimmedQuestion,
        description: description.trim() || undefined,
        options: trimmedOptions,
      });
      onSuccess(poll);
    } catch (err: any) {
      setError(err.message || 'Failed to create poll.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '640px', margin: '40px auto', padding: '0 20px' }}>
      <div className="card" style={{ padding: '32px' }}>
        <div style={{ marginBottom: '24px' }}>
          <span className="brand-badge">NEW POLL BUILDER</span>
          <h2 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '28px', color: 'var(--v-color-deep-ink)', marginTop: '6px' }}>
            Create a New Poll
          </h2>
          <p style={{ color: 'var(--v-color-warm-gray)', fontSize: '14px', marginTop: '4px' }}>
            Fill out the question and options below to launch your poll.
          </p>
        </div>

        {error && (
          <div style={{
            padding: '12px 16px',
            backgroundColor: 'rgba(196, 52, 45, 0.08)',
            borderLeft: '3px solid var(--v-color-error)',
            color: 'var(--v-color-error)',
            fontSize: '13px',
            marginBottom: '20px'
          }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <div>
            <label style={{ display: 'block', fontSize: '13px', fontWeight: 600, marginBottom: '6px' }}>
              Poll Question *
            </label>
            <input
              type="text"
              required
              minLength={5}
              maxLength={300}
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              placeholder="What feature should we build next?"
              style={{
                width: '100%',
                padding: '12px 14px',
                borderRadius: 'var(--v-radius-md)',
                border: '1px solid var(--v-color-stone)',
                fontFamily: 'var(--v-font-body)',
                fontSize: '15px',
                outline: 'none',
              }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '13px', fontWeight: 600, marginBottom: '6px' }}>
              Description / Notes (Optional)
            </label>
            <textarea
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Add background context for your audience..."
              style={{
                width: '100%',
                padding: '10px 14px',
                borderRadius: 'var(--v-radius-md)',
                border: '1px solid var(--v-color-stone)',
                fontFamily: 'var(--v-font-body)',
                fontSize: '14px',
                outline: 'none',
                resize: 'vertical',
              }}
            />
          </div>

          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
              <label style={{ fontSize: '13px', fontWeight: 600 }}>Voting Options (Min 2, Max 10)</label>
              <span style={{ fontSize: '12px', color: 'var(--v-color-warm-gray)' }}>{options.length}/10</span>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
              {options.map((optText, index) => (
                <div key={index} style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                  <span style={{ fontSize: '12px', color: 'var(--v-color-warm-gray)', width: '20px' }}>
                    {index + 1}.
                  </span>
                  <input
                    type="text"
                    required
                    value={optText}
                    onChange={(e) => handleOptionChange(index, e.target.value)}
                    placeholder={`Option ${index + 1}`}
                    style={{
                      flex: 1,
                      padding: '10px 14px',
                      borderRadius: 'var(--v-radius-md)',
                      border: '1px solid var(--v-color-stone)',
                      fontFamily: 'var(--v-font-body)',
                      fontSize: '14px',
                      outline: 'none',
                    }}
                  />
                  {options.length > 2 && (
                    <button
                      type="button"
                      onClick={() => handleRemoveOption(index)}
                      style={{
                        background: 'none',
                        border: 'none',
                        color: 'var(--v-color-error)',
                        fontSize: '18px',
                        cursor: 'pointer',
                        padding: '0 8px',
                      }}
                      title="Remove Option"
                    >
                      ✕
                    </button>
                  )}
                </div>
              ))}
            </div>

            {options.length < 10 && (
              <button
                type="button"
                onClick={handleAddOption}
                style={{
                  marginTop: '12px',
                  background: 'none',
                  border: '1px dashed var(--v-color-stone)',
                  width: '100%',
                  padding: '10px',
                  borderRadius: 'var(--v-radius-md)',
                  color: 'var(--v-color-ink-black)',
                  fontSize: '13px',
                  fontWeight: 600,
                  cursor: 'pointer',
                }}
              >
                + Add Option
              </button>
            )}
          </div>

          <div style={{ display: 'flex', gap: '12px', marginTop: '12px' }}>
            <Button type="submit" variant="primary" size="lg" fullWidth disabled={loading}>
              {loading ? 'Publishing...' : 'Create & Publish Poll'}
            </Button>
            <Button type="button" variant="outline" size="lg" onClick={onCancel}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
