import type { GlobalThemeOverrides } from 'naive-ui'

export const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#2563EB',
    primaryColorHover: '#3B82F6',
    primaryColorPressed: '#1D4ED8',
    primaryColorSuppl: '#2563EB',
    borderRadius: '10px',
    borderRadiusSmall: '8px',
    fontFamily:
      '"Inter", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  },
  Card: {
    borderRadius: '16px',
    borderColor: '#E5EAF3',
    color: '#FFFFFF',
    paddingMedium: '20px',
  },
  Button: {
    borderRadiusMedium: '10px',
    borderRadiusSmall: '8px',
    heightMedium: '36px',
    heightSmall: '32px',
    fontSizeMedium: '14px',
    fontSizeSmall: '13px',
    colorPrimary: '#2563EB',
    colorHoverPrimary: '#3B82F6',
    colorPressedPrimary: '#1D4ED8',
  },
  Input: {
    borderRadius: '10px',
    heightMedium: '38px',
    borderColor: '#E5EAF3',
    borderHover: '#CBD5E1',
    borderFocus: '#2563EB',
  },
  Tag: {
    borderRadius: '999px',
    heightMedium: '26px',
    heightSmall: '22px',
  },
  Dropdown: {
    borderRadius: '12px',
    boxShadow: '0 12px 32px rgba(15, 23, 42, 0.12)',
  },
  Menu: {
    itemHeight: '46px',
    borderRadius: '12px',
    itemTextColorInverted: 'rgba(255, 255, 255, 0.72)',
    itemTextColorHoverInverted: '#FFFFFF',
    itemTextColorActiveInverted: '#FFFFFF',
    itemTextColorChildActiveInverted: 'rgba(255, 255, 255, 0.85)',
    itemIconColorInverted: 'rgba(255, 255, 255, 0.72)',
    itemIconColorHoverInverted: '#FFFFFF',
    itemIconColorActiveInverted: '#FFFFFF',
    itemIconColorChildActiveInverted: 'rgba(255, 255, 255, 0.85)',
    itemColorActiveInverted: 'rgba(37, 99, 235, 0.18)',
    itemColorHoverInverted: 'rgba(255, 255, 255, 0.06)',
    arrowIconColorInverted: 'rgba(255, 255, 255, 0.5)',
  },
  Layout: {
    headerBorderColor: '#E5EAF3',
  },
}
