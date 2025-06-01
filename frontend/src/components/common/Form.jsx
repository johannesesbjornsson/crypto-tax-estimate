import React from 'react';

const Form = ({
  fields,
  formData,
  handleFormChange,
  handleSubmit,
  formErrorMessage,
  successMessage,
  setShowForm,
}) => {
  return (
    <form onSubmit={handleSubmit}>
      {fields.map((field) => {
        if (field.type === 'group') {
          return (
            <div key={field.label} className="setting-row">
              <span className="setting-label">{field.label}:</span>
              {field.fields.map((subField, i) => (
                <React.Fragment key={subField.name}>
                  {subField.type === 'select' ? (
                    <select
                      className="setting-value"
                      name={subField.name}
                      value={formData[subField.name]}
                      onChange={handleFormChange}
                    >
                      {subField.options.map((option) => (
                        <option key={option} value={option}>{option}</option>
                      ))}
                    </select>
                  ) : (
                    <input
                      className="setting-value"
                      name={subField.name}
                      type={subField.type}
                      value={formData[subField.name]}
                      onChange={handleFormChange}
                    />
                  )}
                  {i < field.fields.length - 1 && <span className="date-separator">/</span>}
                </React.Fragment>
              ))}
            </div>
          );
        }

        return (
          <div key={field.name} className="setting-row">
            <span className="setting-label">{field.label}:</span>
            {field.type === 'display' ? (
              <span className="setting-value">{formData[field.name] || '\u00A0'}</span>
            ) : field.type === 'select' ? (
              <select
                className="setting-value"
                name={field.name}
                value={formData[field.name]}
                onChange={handleFormChange}
                required={field.required}
              >
                {field.options.map((option) => (
                  <option key={option} value={option}>{option}</option>
                ))}
              </select>
            ) : (
              <input
                className="setting-value"
                name={field.name}
                type={field.type}
                step={field.step}
                accept={field.accept}
                value={field.type === 'file' ? undefined : formData[field.name]}
                onChange={handleFormChange}
                required={field.required}
              />
            )}
          </div>
        );
      })}

      <div className="save-row">
        {successMessage && (
          <span className="save-message success">{successMessage}</span>
        )}
        {formErrorMessage && (
          <span className="save-message error">{formErrorMessage}</span>
        )}
        <button type="submit">Save</button>
        {typeof setShowForm === 'function' && (
          <button type="button" onClick={() => setShowForm(false)}>Cancel</button>
        )}
      </div>
    </form>
  );
};

export default Form;