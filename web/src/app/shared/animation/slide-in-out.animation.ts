import { animate, state, style, transition, trigger } from '@angular/animations';

/**
 * A slide in-out animation
 */
export const slideInOut = trigger('slideInOut', [
  state('true', style({ transform: 'translateX(0)' })),
  state('false', style({ transform: 'translateX(-25.6rem)' })),
  transition('* => true', [
    animate('.4s ease-in-out', style({ transform: 'translateX(0)' })),
  ]),
  transition('true => false', [
    animate('.4s ease-in-out', style({ transform: 'translateX(-25.6rem)' })),
  ]),
]);
