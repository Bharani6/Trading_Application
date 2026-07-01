/**
 * Formats a number to 2 decimal places with comma separation
 * @param {Number|String} val - The value to format
 * @returns {String} Formatted number
 */
export const formatNumber = (val) => {
  if (!val && val !== 0) return '0.00'
  return Number(val).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

/**
 * Formats a date string to a readable format
 * @param {String|Date} dateStr - The date to format
 * @returns {String} Formatted date (e.g., "Jan 1, 10:30 AM")
 */
export const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}
