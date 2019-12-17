import { AfterViewInit, Directive, ElementRef, HostListener } from '@angular/core';

@Directive({ selector: '[appHeightAsScroll]' })
export class HeightAsScrollDirective implements AfterViewInit {

  constructor(private host: ElementRef) { }

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
    const elem: HTMLElement = this.host.nativeElement as HTMLElement;

    elem.style.height = elem.scrollHeight.toString() + 'px';
  }

}
