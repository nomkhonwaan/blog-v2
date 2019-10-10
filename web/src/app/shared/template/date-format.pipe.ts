import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'dateFormat' })
export class DateFormatPipe implements PipeTransform {

  transform(value?: string): string {
    let d: Date;

    try {
      d = value ? new Date(value) : new Date();
    } catch (err) {
      console.error(`an error has occurred while parsing date string ${value}: ${err}`);
      d = new Date();
    }

    return `${this.getFullMonth(d.getMonth())} ${d.getDate()}, ${d.getFullYear()}`;
  }

  getFullMonth(month: number): string {
    const fullMonthNames: string[] = [
      'January',
      'Febuary',
      'March',
      'April',
      'May',
      'June',
      'July',
      'August',
      'September',
      'October',
      'November',
      'December',
    ];

    return fullMonthNames[month];
  }

}
