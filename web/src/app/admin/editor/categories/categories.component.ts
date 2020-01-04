import { Component, OnInit } from '@angular/core';
import { faSearch, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { AbstractPostEditorComponent } from '../abstract-post-editor.component';

type Archive = Category | Tag;

@Component({
  selector: 'app-post-categories-editor',
  templateUrl: './categories.component.html',
  styleUrls: ['./categories.component.scss'],
})
export class PostCategoriesEditorComponent extends AbstractPostEditorComponent implements OnInit {

  /**
   * List of categories or tags to-be rendered as sidebar menu-item(s)
   */
  archives: Array<Archive>;

  /**
   * List of categories or tags are matched with the search keyword
   */
  matchedSearchArchives: Array<Archive> = [];

  /**
   * List of categories or tags are shown in the dropdown menu
   */
  dropdownItems: Array<DropdownItem> = [];

  /**
   * List of categories or tags are selected from the dropdown menu
   */
  selectedArchives: Array<Archive> = [];

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faSearch,
  };

  ngOnInit(): void {
    this.getAll('categories');
  }

  onChange(selectedItem: DropdownItem): void {
    this.selectedArchives = this.selectedArchives
      .concat(this.archives.find(({ slug }: Archive): boolean => selectedItem.value === slug))
      .filter((val, i, self) => self.indexOf(val) === i);
    // this.mutate(
    //   `
    //       mutation {
    //         updatePostStatus(slug: $slug, status: $status) {
    //           ...EditablePost
    //         }
    //       }
    //     `,
    //   {
    //     slug: this.post.slug,
    //     status: selectedItem.value.toString(),
    //   },
    // ).subscribe((post: Post): void => {
    //   this.changeSuccess.emit(post);

    //   this.currentStatus = update(selectedItem, { label: { $set: selectedItem.label.toUpperCase() } });
    // });
  }

  onBlur(): void {
    setTimeout(() => { this.dropdownItems = []; }, 800);
  }

  onKeyPress(keyword: string): void {
    this.findAllArchives(keyword);
  }

  protected getAll(type: string): void {
    this.apollo.query({
      query: gql`
        {
          archives: ${type} { name slug }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ archives: Array<Archive> }>): Array<Archive> => result.data.archives),
    ).subscribe((archives: Array<Archive>): void => {
      this.archives = archives;
    });
  }

  private findAllArchives(keyword: string): void {
    if (this.archives.length === 0) {
      return;
    }

    const re: RegExp = new RegExp(keyword, 'i');

    this.matchedSearchArchives =
      this.archives.reduce((result: Array<Archive>, archive: Archive): Array<Archive> => {
        if (re.test(archive.name)) {
          return result.concat(archive);
        }

        return result;
      }, []);

    this.dropdownItems =
      this.matchedSearchArchives.map(({ name, slug }: Archive): DropdownItem => ({ label: name, value: slug }));
  }
}
