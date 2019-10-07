import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'split',
})
export class StringSplit implements PipeTransform {

  transform(value: string, separator: string | RegExp, limit?: number): string[] {
    return value.split(separator, limit);
  }

}
