import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'split' })
export class StringSplitPipe implements PipeTransform {

  transform(value: string, separator: string | RegExp, limit?: number): string[] {
    return value.split(separator, limit);
  }

}
