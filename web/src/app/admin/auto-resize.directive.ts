import { Directive, AfterViewInit, ElementRef, HostListener } from '@angular/core';

@Directive({ selector: '[appAutoResize]' })
export class AutoResizeDirective implements AfterViewInit {

  constructor(private elementRef: ElementRef) { }

  ngAfterViewInit(): void {
    setTimeout(() => this.resize());
  }

  @HostListener('change')
  onChange(): void {
    this.resize();
  }

  @HostListener('document:keypress', ['$event'])
  onKeyPress(event: KeyboardEvent): void {
    this.resize();
  }

  private resize(): void {
    const elem: HTMLElement = this.elementRef.nativeElement as HTMLElement;

    elem.style.height = elem.scrollHeight.toString() + 'px';
  }

}
